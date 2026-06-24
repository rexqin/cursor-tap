import { act, renderHook, waitFor } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { useRecords } from './use-records';
import { sampleGrpcRecord, sampleRecord } from '@/test/fixtures/records';

describe('useRecords', () => {
  it('fetches initial records from API via MSW', async () => {
    const { result } = renderHook(() => useRecords());

    await act(async () => {
      await result.current.fetchInitialRecords();
    });

    await waitFor(() => {
      expect(result.current.allRecords.length).toBeGreaterThan(0);
    });
    expect(result.current.allSessions.length).toBeGreaterThan(0);
  });

  it('deduplicates records added via WebSocket', async () => {
    const { result } = renderHook(() => useRecords());

    act(() => {
      result.current.addRecord(sampleRecord);
      result.current.addRecord(sampleRecord);
    });

    expect(result.current.allRecords).toHaveLength(1);
  });

  it('does not add records when paused', () => {
    const { result } = renderHook(() => useRecords());

    act(() => {
      result.current.setIsPaused(true);
    });
    act(() => {
      result.current.addRecord(sampleGrpcRecord);
    });

    expect(result.current.allRecords).toHaveLength(0);
  });

  it('filters records by search query', async () => {
    const { result } = renderHook(() => useRecords());

    await act(async () => {
      await result.current.fetchInitialRecords();
    });

    act(() => {
      result.current.setSearchQuery('AgentService');
    });

    await waitFor(() => {
      expect(result.current.records.every((r) =>
        r.url?.includes('AgentService') ||
        r.grpc_service?.includes('AgentService'),
      )).toBe(true);
    });
  });

  it('clearRecords resets state', async () => {
    const { result } = renderHook(() => useRecords());

    await act(async () => {
      await result.current.fetchInitialRecords();
    });

    act(() => {
      result.current.clearRecords();
    });

    expect(result.current.allRecords).toHaveLength(0);
    expect(result.current.allSessions).toHaveLength(0);
  });
});
