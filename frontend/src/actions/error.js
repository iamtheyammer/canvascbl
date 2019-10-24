// 502 BAD GATEWAY from proxy
export const CANVAS_PROXY_ERROR = 'CANVAS_PROXY_ERROR';

export const CHECKOUT_ERROR = 'CHECKOUT_ERROR';

export const PLUS_ERROR = 'PLUS_ERROR';

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
  };
}

export function plusError(id, res) {
  return {
    type: PLUS_ERROR,
    id,
    res
  };
}
