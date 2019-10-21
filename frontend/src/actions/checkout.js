import makeCheckoutRequest from '../util/checkout/makeCheckoutRequest';
import { checkoutError } from './error';
import { startLoading, endLoading } from './loading';

export const CHECKOUT_GOT_PRODUCTS = 'CHECKOUT_GOT_PRODUCTS';

export const CHECKOUT_GOT_CHECKOUT_SESSION = 'CHECKOUT_GOT_CHECKOUT_SESSION';

function gotProducts(products) {
  return {
    type: CHECKOUT_GOT_PRODUCTS,
    products
  };
}

export function getProducts(id) {
  return async dispatch => {
    dispatch(startLoading(id));
    try {
      const products = await makeCheckoutRequest('products');
      dispatch(gotProducts(products.data));
    } catch (e) {
      dispatch(checkoutError(id, e.response));
    }
    dispatch(endLoading(id));
  };
}

function gotCheckoutSession(session) {
  return {
    type: CHECKOUT_GOT_CHECKOUT_SESSION,
    session
  };
}

export function getCheckoutSession(id, productId, email) {
  return async dispatch => {
    dispatch(startLoading(id));
    try {
      const session = await makeCheckoutRequest('session', {
        productId,
        email
      });
      dispatch(gotCheckoutSession(session.data));
    } catch (e) {
      dispatch(checkoutError(id, e.response));
    }
    dispatch(endLoading(id));
  };
}
