import { http, HttpResponse } from 'msw';
import { sampleRecords } from '../fixtures/records';

export const handlers = [
  http.get('http://localhost:9090/api/records', () => {
    return HttpResponse.json(sampleRecords);
  }),
];
