'use client';

import { useState, useCallback, useMemo, useRef } from 'react';
import { Record } from '@/lib/types';
import {
  MAX_RECORDS,
  applyGrpcFromRecord,
  computeAvailableFilters,
  computeMethodCounts,
  computeSessions,
  filterRecords,
  filterSessions,
  mergeRecordsDedup,
  recordKey,
  trimRecordsToMax,
} from '@/lib/record-utils';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090';

export function useRecords() {
  const [records, setRecords] = useState<Record[]>([]);
  const [selectedSession, setSelectedSession] = useState<string | null>(null);
  const [selectedRecordKey, setSelectedRecordKey] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isPaused, setIsPaused] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [initialized, setInitialized] = useState(false);
  const [selectedServices, setSelectedServices] = useState<Set<string>>(new Set());
  const [selectedMethods, setSelectedMethods] = useState<Set<string>>(new Set());

  const recordsMap = useRef(new Map<string, Record>());

  const addRecord = useCallback(
    (record: Record) => {
      if (isPaused) return;

      const key = recordKey(record);
      recordsMap.current.set(key, record);

      setRecords((prev) => {
        const exists = prev.some((r) => recordKey(r) === key);
        if (exists) {
          return prev;
        }

        const newRecords = [...prev, record];
        if (newRecords.length > MAX_RECORDS) {
          const trimmed = trimRecordsToMax(newRecords, MAX_RECORDS);
          const removed = newRecords.slice(0, newRecords.length - MAX_RECORDS);
          for (const r of removed) {
            recordsMap.current.delete(recordKey(r));
          }
          return trimmed;
        }
        return newRecords;
      });
    },
    [isPaused],
  );

  const selectedRecord = useMemo(() => {
    if (!selectedRecordKey) return null;
    return recordsMap.current.get(selectedRecordKey) || null;
  }, [selectedRecordKey]);

  const setSelectedRecord = useCallback((record: Record | null) => {
    if (!record) {
      setSelectedRecordKey(null);
    } else {
      const key = recordKey(record);
      recordsMap.current.set(key, record);
      setSelectedRecordKey(key);
    }
  }, []);

  const sessions = useMemo(() => computeSessions(records), [records]);
  const availableFilters = useMemo(() => computeAvailableFilters(sessions), [sessions]);
  const methodCounts = useMemo(() => computeMethodCounts(sessions), [sessions]);
  const filteredSessions = useMemo(
    () => filterSessions(sessions, selectedServices, selectedMethods),
    [sessions, selectedServices, selectedMethods],
  );
  const filteredRecords = useMemo(
    () => filterRecords(records, selectedSession, searchQuery),
    [records, selectedSession, searchQuery],
  );

  const fetchAndMergeRecords = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/api/records?limit=100`);
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        for (const r of data as Record[]) {
          recordsMap.current.set(recordKey(r), r);
        }

        setRecords((prev) => mergeRecordsDedup(prev, data as Record[]));
      }
    } catch (e) {
      console.error('Failed to fetch records:', e);
    }
  }, []);

  const fetchInitialRecords = useCallback(async () => {
    if (initialized) return;
    await fetchAndMergeRecords();
    setInitialized(true);
  }, [initialized, fetchAndMergeRecords]);

  const recoverData = useCallback(async () => {
    await fetchAndMergeRecords();
  }, [fetchAndMergeRecords]);

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
