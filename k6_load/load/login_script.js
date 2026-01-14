import { check, sleep } from 'k6';
import api from 'k6/x/load';

export const options = {
  stages: [
    { duration: '1m', target: 100 },   // 1分钟升到100并发
    { duration: '5m', target: 100 },   // 稳定100并发5分钟
    { duration: '1m', target: 500 },   // 再冲到500并发
    { duration: '3m', target: 500 },
    { duration: '2m', target: 0 },     // 降下来
  ],
  thresholds: {
    'checks': ['rate>0.99'],                    // 总体检查通过率 >99%
    'login_success': ['rate>0.99'],             // 自定义指标
    'http_req_duration': ['p(95)<1000'],        // 95分位 <1s
  },
};

export default function () {
  const userNo = __VU; // 每个虚拟用户一个唯一用户名
  const resp = api.loginY1(`testuser${userNo}`, 'password123');

  // 你的三个断言要求
  const success = check(resp, {
    'code == 0':     (r) => r.code === 0,
    'msg == Succeed': (r) => r.msg === 'Succeed',
    'data 不为空':    (r) => r.data !== null && r.data !== undefined && Object.keys(r.data).length > 0,
  });

  // 自定义指标，报表里单独展示登录成功率
  check(success, { 'login_success': (ok) => ok });

  sleep(0.5 + Math.random() * 1); // 模拟真实用户思考时间
}