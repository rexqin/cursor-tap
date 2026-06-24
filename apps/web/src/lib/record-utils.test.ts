import { describe, expect, it } from 'vitest';
import {
  PREVIEW_MAX_CHARS,
  applyGrpcFromRecord,
  computeAvailableFilters,
  computeMethodCounts,
  computeSessions,
  filterRecords,
  filterSessions,
  mergeRecordsDedup,
  parseGrpcFromUrl,
  recordKey,
  trimRecordsToMax,
  truncatePreview,
} from './record-utils';
import { Record, SessionInfo } from './types';
import { sampleGrpcRecord, sampleRecord } from '@/test/fixtures/records';

describe('parseGrpcFromUrl', () => {
  it.each([
    ['/aiserver.v1.ChatService/StreamChat', { service: 'aiserver.v1.ChatService', method: 'StreamChat' }],
    ['aiserver.v1.AgentService/Run', { service: 'aiserver.v1.AgentService', method: 'Run' }],
  ])('parses %s', (url, expected) => {
    expect(parseGrpcFromUrl(url)).toEqual(expected);
  });

  it.each([undefined, '', '/no-method', 'single-segment'])('returns null for %s', (url) => {
    expect(parseGrpcFromUrl(url)).toBeNull();
  });
});

describe('truncatePreview', () => {
  it('returns undefined for empty input', () => {
    expect(truncatePreview(undefined)).toBeUndefined();
  });

  it('returns short strings unchanged', () => {
    expect(truncatePreview('hello')).toBe('hello');
  });

  it('truncates long strings to PREVIEW_MAX_CHARS', () => {
    const long = 'x'.repeat(PREVIEW_MAX_CHARS + 50);
    expect(truncatePreview(long)).toHaveLength(PREVIEW_MAX_CHARS);
  });
});

describe('recordKey', () => {
  it('combines session and index', () => {
    expect(recordKey({ session: 'abc', index: 3 } as Record)).toBe('abc-3');
  });
});

describe('applyGrpcFromRecord', () => {
  it('sets grpc fields from record.grpc_service', () => {
    const session: SessionInfo = {
      id: 's1',
      seq: 1,
      host: '',
      record_count: 1,
      first_ts: '',
      last_ts: '',
      request_size: 0,
      response_size: 0,
    };
    applyGrpcFromRecord(session, sampleGrpcRecord);
    expect(session.grpc_service).toBe('aiserver.v1.ChatService');
    expect(session.grpc_method).toBe('StreamChat');
  });

  it('parses grpc info from request URL when grpc fields absent', () => {
    const session: SessionInfo = {
      id: 's1',
      seq: 1,
      host: '',
      record_count: 1,
      first_ts: '',
      last_ts: '',
      request_size: 0,
      response_size: 0,
    };
    applyGrpcFromRecord(session, sampleRecord);
    expect(session.grpc_service).toBe('aiserver.v1.ChatService');
    expect(session.grpc_method).toBe('StreamChat');
  });

  it('does not overwrite existing grpc_service', () => {
    const session: SessionInfo = {
      id: 's1',
      seq: 1,
      host: '',
      record_count: 1,
      first_ts: '',
      last_ts: '',
      request_size: 0,
      response_size: 0,
      grpc_service: 'existing.Service',
      grpc_method: 'Existing',
    };
    applyGrpcFromRecord(session, sampleRecord);
    expect(session.grpc_service).toBe('existing.Service');
  });
});

describe('computeSessions', () => {
  it('aggregates records into sessions sorted by seq desc', () => {
    const sessions = computeSessions([sampleRecord, sampleGrpcRecord]);
    expect(sessions).toHaveLength(1);
    expect(sessions[0].record_count).toBe(2);
    expect(sessions[0].grpc_service).toBe('aiserver.v1.ChatService');
  });

  it('accumulates request/response sizes', () => {
    const sessions = computeSessions([sampleRecord, sampleGrpcRecord]);
    expect(sessions[0].request_size).toBe(192);
  });
});

describe('filterSessions', () => {
  const sessions = computeSessions([
    sampleRecord,
    {
      ...sampleRecord,
      session: 'sess-2',
      seq: 2,
      url: '/aiserver.v1.AgentService/Run',
    },
  ]);

  it('returns all when no filters selected', () => {
    expect(filterSessions(sessions, new Set(), new Set())).toHaveLength(2);
  });

  it('filters by service', () => {
    const filtered = filterSessions(sessions, new Set(['aiserver.v1.ChatService']), new Set());
    expect(filtered).toHaveLength(1);
    expect(filtered[0].id).toBe('sess-1');
  });

  it('filters by full method name', () => {
    const filtered = filterSessions(
      sessions,
      new Set(),
      new Set(['aiserver.v1.AgentService.Run']),
    );
    expect(filtered).toHaveLength(1);
    expect(filtered[0].id).toBe('sess-2');
  });
});

describe('filterRecords', () => {
  const records = [sampleRecord, sampleGrpcRecord];

  it('filters by session', () => {
    expect(filterRecords(records, 'sess-1', '')).toHaveLength(2);
    expect(filterRecords(records, 'other', '')).toHaveLength(0);
  });

  it('filters by search query (case insensitive)', () => {
    expect(filterRecords(records, null, 'streamchat')).toHaveLength(2);
    expect(filterRecords(records, null, 'nomatch')).toHaveLength(0);
  });
});

describe('mergeRecordsDedup', () => {
  it('merges new records and skips duplicates', () => {
    const merged = mergeRecordsDedup([sampleRecord], [sampleRecord, sampleGrpcRecord]);
    expect(merged).toHaveLength(2);
  });

  it('sorts by seq then index', () => {
    const later: Record = { ...sampleGrpcRecord, index: 5, seq: 3 };
    const merged = mergeRecordsDedup([later], [sampleRecord]);
    expect(merged[0].seq).toBe(1);
    expect(merged[1].seq).toBe(3);
  });
});

describe('trimRecordsToMax', () => {
  it('keeps the newest records when over limit', () => {
    const records = Array.from({ length: 5 }, (_, i) => ({
      ...sampleRecord,
      index: i,
    }));
    const trimmed = trimRecordsToMax(records, 3);
    expect(trimmed).toHaveLength(3);
    expect(trimmed[0].index).toBe(2);
  });
});

describe('computeAvailableFilters and computeMethodCounts', () => {
  it('builds service/method map and counts', () => {
    const sessions = computeSessions([
      sampleRecord,
      { ...sampleRecord, session: 'sess-2', seq: 2, url: '/aiserver.v1.AgentService/Run' },
    ]);
    const filters = computeAvailableFilters(sessions);
    expect(filters.get('aiserver.v1.ChatService')?.has('StreamChat')).toBe(true);
    expect(filters.get('aiserver.v1.AgentService')?.has('Run')).toBe(true);

    const counts = computeMethodCounts(sessions);
    expect(counts.get('aiserver.v1.ChatService.StreamChat')).toBe(1);
  });
});
