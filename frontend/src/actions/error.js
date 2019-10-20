// 502 BAD GATEWAY from proxy
export const CANVAS_PROXY_ERROR = 'CANVAS_PROXY_ERROR';

export const CHECKOUT_ERROR = 'CHECKOUT_ERROR';

export function canvasProxyError(id, res) {
  return {
    type: CANVAS_PROXY_ERROR,
    id,
    res
  };
}

export function checkoutError(id, res) {
  return {
    type: CHECKOUT_ERROR,
    id,
    res
  }
}