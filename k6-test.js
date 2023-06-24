import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { target: 200, duration: '15s' },
    { target: 0, duration: '15s' },
  ],
};

export default function () {
  const result = http.get('https://example.local:9443/');
  check(result, {
    'http response status code is 200': result.status === 200,
  });
}