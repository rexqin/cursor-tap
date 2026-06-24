'use client';

import { useState, useCallback, useMemo, useRef } from 'react';
import { Record, SessionInfo } from '@/lib/types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090';
const MAX_RECORDS = 2000;
const PREVIEW_MAX_CHARS = 200;

function truncatePreview(data?: string): string | undefined {
  if (!data) return undefined;
  return data.length > PREVIEW_MAX_CHARS ? data.slice(0, PREVIEW_MAX_CHARS) : data;
}

/** Parse /aiserver.v1.Foo/Bar → { service, method } */
function parseGrpcFromUrl(url?: string): { service: string; method: string } | null {
  if (!url) return null;
  const path = url.startsWith('/') ? url.slice(1) : url;
  const idx = path.lastIndexOf('/');
  if (idx <= 0) return null;
  return { service: path.slice(0, idx), method: path.slice(idx + 1) };
}

function applyGrpcFromRecord(session: SessionInfo, record: Record) {
  if (record.grpc_service) {
    if (!session.grpc_service) {
      session.grpc_service = record.grpc_service;
      session.grpc_method = record.grpc_method;
    }
    return;
  }
  if (record.type === 'request' && record.url && !session.grpc_service) {
    const parsed = parseGrpcFromUrl(record.url);
    if (parsed) {
      session.grpc_service = parsed.service;
      session.grpc_method = parsed.method;
    }
  }
}

export function useRecords() {
  const [records, setRecords] = useState<Record[]>([]);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);
  // Store record key instead of object to keep reference stable
  const [selectedRecordKey, setSelectedRecordKey] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isPaused, setIsPaused] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [initialized, setInitialized] = useState(false);
  // Filter state - multi-select
  const [selectedServices, setSelectedServices] = useState<Set<string>>(new Set());
  const [selectedMethods, setSelectedMethods] = useState<Set<string>>(new Set());

  // Cache for quick record lookup
  const recordsMap = useRef(new Map<string, Record>());

  // Add a new record from WebSocket (with deduplication)
  const addRecord = useCallback((record: Record) => {
    if (isPaused) return;

    const key = `${record.session}-${record.index}`;

    // Update cache
    recordsMap.current.set(key, record);

    setRecords((prev) => {
      // Deduplicate by session + index
      const exists = prev.some((r) => `${r.session}-${r.index}` === key);
      if (exists) {
        return prev;
      }

      const newRecords = [...prev, record];
      // Keep only last MAX_RECORDS in browser
      if (newRecords.length > MAX_RECORDS) {
        // Clean up old entries from cache
        const removed = newRecords.slice(0, newRecords.length - MAX_RECORDS);
        for (const r of removed) {
          recordsMap.current.delete(`${r.session}-${r.index}`);
        }
        return newRecords.slice(-MAX_RECORDS);
      }
      return newRecords;
    });
  }, [isPaused]);

  // Compute selected record from key (stable reference)
  const selectedRecord = useMemo(() => {
    if (!selectedRecordKey) return null;
    return recordsMap.current.get(selectedRecordKey) || null;
  }, [selectedRecordKey]);

  // Wrapper to set record by object (finds key)
  const setSelectedRecord = useCallback((record: Record | null) => {
    if (!record) {
      setSelectedRecordKey(null);
    } else {
      const key = `${record.session}-${record.index}`;
      // Ensure it's in cache
      recordsMap.current.set(key, record);
      setSelectedRecordKey(key);
    }
  }, []);

  // Compute sessions (RPC calls) from records (browser-side)
  const sessions = useMemo(() => {
    const sessionMap = new Map<string, SessionInfo>();

    for (const record of records) {
      const existing = sessionMap.get(record.session);
      if (existing) {
        existing.record_count++;
        if (record.ts > existing.last_ts) {
          existing.last_ts = record.ts;
        }
        if (record.ts < existing.first_ts) {
          existing.first_ts = record.ts;
        }
        // Update gRPC info from grpc records or request URL
        applyGrpcFromRecord(existing, record);
        // Update URL from request
        if (record.type === 'request' && record.url && !existing.url) {
          existing.url = record.url;
        }
        // Update sizes
        if (record.direction === 'C2S' && record.size) {
          existing.request_size += record.size;
        }
        if (record.direction === 'S2C' && record.size) {
          existing.response_size += record.size;
        }
        // Capture first C2S gRPC data as preview
        if (record.type === 'grpc' && record.direction === 'C2S' && record.grpc_data && !existing.grpc_preview) {
          existing.grpc_preview = truncatePreview(record.grpc_data);
        }
      } else {
        const session: SessionInfo = {
          id: record.session,
          seq: record.seq,
          host: record.host || '',
          record_count: 1,
          first_ts: record.ts,
          last_ts: record.ts,
          url: record.type === 'request' ? record.url : undefined,
          request_size: record.direction === 'C2S' ? (record.size || 0) : 0,
          response_size: record.direction === 'S2C' ? (record.size || 0) : 0,
          grpc_preview: record.type === 'grpc' && record.direction === 'C2S'
            ? truncatePreview(record.grpc_data)
            : undefined,
        };
        applyGrpcFromRecord(session, record);
        sessionMap.set(record.session, session);
      }
    }

    // Sort by seq descending (newest first)
    return Array.from(sessionMap.values()).sort((a, b) => b.seq - a.seq);
  }, [records]);

  // Extract available services and methods for filter
  const availableFilters = useMemo(() => {
    const services = new Map<string, Set<string>>();

    for (const session of sessions) {
      if (session.grpc_service) {
        if (!services.has(session.grpc_service)) {
          services.set(session.grpc_service, new Set());
        }
        if (session.grpc_method) {
          services.get(session.grpc_service)!.add(session.grpc_method);
        }
      }
    }

    return services;
  }, [sessions]);

  // Count sessions per method
  const methodCounts = useMemo(() => {
    const counts = new Map<string, number>();

    for (const session of sessions) {
      if (session.grpc_service && session.grpc_method) {
        const key = `${session.grpc_service}.${session.grpc_method}`;
        counts.set(key, (counts.get(key) || 0) + 1);
      }
    }

    return counts;
  }, [sessions]);

  // Filtered sessions by service/method selection
  const filteredSessions = useMemo(() => {
    if (selectedServices.size === 0 && selectedMethods.size === 0) {
      return sessions;
    }

    return sessions.filter((s) => {
      if (!s.grpc_service) return false;

      // If specific methods selected, check service.method
      if (selectedMethods.size > 0) {
        const fullMethod = `${s.grpc_service}.${s.grpc_method}`;
        return selectedMethods.has(fullMethod);
      }

      // Otherwise check service
      if (selectedServices.size > 0) {
        return selectedServices.has(s.grpc_service);
      }

      return true;
    });
  }, [sessions, selectedServices, selectedMethods]);

  // Filter records by selected session and search query
  const filteredRecords = useMemo(() => {
    let result = records;

    if (selectedSession) {
      result = result.filter((r) => r.session === selectedSession);
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      result = result.filter((r) => {
        return (
          r.url?.toLowerCase().includes(query) ||
          r.grpc_service?.toLowerCase().includes(query) ||
          r.grpc_method?.toLowerCase().includes(query) ||
          r.grpc_data?.toLowerCase().includes(query) ||
          r.body?.toLowerCase().includes(query) ||
          r.host?.toLowerCase().includes(query)
        );
      });
    }

    return result;
  }, [records, selectedSession, searchQuery]);

  // Fetch and merge records from API (with deduplication)
  const fetchAndMergeRecords = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/api/records?limit=100`);
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        // Update cache
        for (const r of data as Record[]) {
          recordsMap.current.set(`${r.session}-${r.index}`, r);
        }

        setRecords((prev) => {
          // Merge and deduplicate
          const existingKeys = new Set(prev.map((r) => `${r.session}-${r.index}`));
          const newRecords = (data as Record[]).filter(
            (r) => !existingKeys.has(`${r.session}-${r.index}`)
          );
          if (newRecords.length === 0) return prev;
          return [...prev, ...newRecords].sort((a, b) => {
            if (a.seq !== b.seq) return a.seq - b.seq;
            return a.index - b.index;
          });
        });
      }
    } catch (e) {
      console.error('Failed to fetch records:', e);
    }
  }, []);

  // Fetch initial records from API (only once)
  const fetchInitialRecords = useCallback(async () => {
    if (initialized) return;
    await fetchAndMergeRecords();
    setInitialized(true);
  }, [initialized, fetchAndMergeRecords]);

  // Recover data on reconnect
  const recoverData = useCallback(async () => {
    console.log('Recovering data after reconnect...');
    await fetchAndMergeRecords();
  }, [fetchAndMergeRecords]);

  // Clear all records
  const clearRecords = useCallback(() => {
    setRecords([]);
    setSelectedSession(null);
    setSelectedRecordKey(null);
    recordsMap.current.clear();
  }, []);

  return {
    records: filteredRecords,
    allRecords: records,
    sessions: filteredSessions,
    allSessions: sessions,
    availableFilters,
    methodCounts,
    selectedSession,
    selectedRecord,
    selectedServices,
    selectedMethods,
    isConnected,
    isPaused,
    searchQuery,
    setSelectedSession,
    setSelectedRecord,
    setSelectedServices,
    setSelectedMethods,
    setIsConnected,
    setIsPaused,
    setSearchQuery,
    addRecord,
    fetchInitialRecords,
    recoverData,
    clearRecords,
  };
}
