import { Record, SessionInfo } from './types';

export const MAX_RECORDS = 2000;
export const PREVIEW_MAX_CHARS = 200;

export function recordKey(record: Pick<Record, 'session' | 'index'>): string {
  return `${record.session}-${record.index}`;
}

export function truncatePreview(data?: string): string | undefined {
  if (!data) return undefined;
  return data.length > PREVIEW_MAX_CHARS ? data.slice(0, PREVIEW_MAX_CHARS) : data;
}

/** Parse /aiserver.v1.Foo/Bar → { service, method } */
export function parseGrpcFromUrl(url?: string): { service: string; method: string } | null {
  if (!url) return null;
  const path = url.startsWith('/') ? url.slice(1) : url;
  const idx = path.lastIndexOf('/');
  if (idx <= 0) return null;
  return { service: path.slice(0, idx), method: path.slice(idx + 1) };
}

export function applyGrpcFromRecord(session: SessionInfo, record: Record): void {
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

export function computeSessions(records: Record[]): SessionInfo[] {
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
      applyGrpcFromRecord(existing, record);
      if (record.type === 'request' && record.url && !existing.url) {
        existing.url = record.url;
      }
      if (record.direction === 'C2S' && record.size) {
        existing.request_size += record.size;
      }
      if (record.direction === 'S2C' && record.size) {
        existing.response_size += record.size;
      }
      if (
        record.type === 'grpc' &&
        record.direction === 'C2S' &&
        record.grpc_data &&
        !existing.grpc_preview
      ) {
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
        grpc_preview:
          record.type === 'grpc' && record.direction === 'C2S'
            ? truncatePreview(record.grpc_data)
            : undefined,
      };
      applyGrpcFromRecord(session, record);
      sessionMap.set(record.session, session);
    }
  }

  return Array.from(sessionMap.values()).sort((a, b) => b.seq - a.seq);
}

export function computeAvailableFilters(
  sessions: SessionInfo[],
): Map<string, Set<string>> {
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
}

export function computeMethodCounts(sessions: SessionInfo[]): Map<string, number> {
  const counts = new Map<string, number>();

  for (const session of sessions) {
    if (session.grpc_service && session.grpc_method) {
      const key = `${session.grpc_service}.${session.grpc_method}`;
      counts.set(key, (counts.get(key) || 0) + 1);
    }
  }

  return counts;
}

export function filterSessions(
  sessions: SessionInfo[],
  selectedServices: Set<string>,
  selectedMethods: Set<string>,
): SessionInfo[] {
  if (selectedServices.size === 0 && selectedMethods.size === 0) {
    return sessions;
  }

  return sessions.filter((s) => {
    if (!s.grpc_service) return false;

    if (selectedMethods.size > 0) {
      const fullMethod = `${s.grpc_service}.${s.grpc_method}`;
      return selectedMethods.has(fullMethod);
    }

    if (selectedServices.size > 0) {
      return selectedServices.has(s.grpc_service);
    }

    return true;
  });
}

export function filterRecords(
  records: Record[],
  selectedSession: string | null,
  searchQuery: string,
): Record[] {
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
}

export function mergeRecordsDedup(prev: Record[], incoming: Record[]): Record[] {
  const existingKeys = new Set(prev.map(recordKey));
  const newRecords = incoming.filter((r) => !existingKeys.has(recordKey(r)));
  if (newRecords.length === 0) return prev;
  return [...prev, ...newRecords].sort((a, b) => {
    if (a.seq !== b.seq) return a.seq - b.seq;
    return a.index - b.index;
  });
}

export function trimRecordsToMax(records: Record[], max: number): Record[] {
  if (records.length <= max) return records;
  return records.slice(-max);
}
