import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { target: 1000, duration: '30s' },
    { target: 0, duration: '30s' },
  ],
};

export default function () {
  const result = http.get('http://localhost:8081/');
  check(result, {
    'http response status code is 200': result.status === 200,
  });
}