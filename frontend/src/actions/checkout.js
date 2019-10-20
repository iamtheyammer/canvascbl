import makeCheckoutRequest from "../util/checkout/makeCheckoutRequest";

import { checkoutError } from './error';

import { startLoading, endLoading } from './loading';

export const CHECKOUT_GOT_CHECKOUT_TOKEN = 'CHECKOUT_GOT_CHECKOUT_TOKEN';

function gotCheckoutSession(checkoutSession) {
  return {
    type: CHECKOUT_GOT_CHECKOUT_TOKEN,
    checkoutSession
  }
}

export function getCheckoutSession(id) {
  return async dispatch => {
    dispatch(startLoading(id));
    try {
      const session = await makeCheckoutRequest('session');
      dispatch(gotCheckoutSession(session.data))
    } catch (e) {
      dispatch(checkoutError(id, e.response))
    }
    dispatch(endLoading(id));
  }
}
