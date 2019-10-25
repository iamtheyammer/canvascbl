import makeCheckoutRequest from '../util/checkout/makeCheckoutRequest';
import { checkoutError } from './error';
import { startLoading, endLoading } from './loading';

export const CHECKOUT_GOT_PRODUCTS = 'CHECKOUT_GOT_PRODUCTS';

export const CHECKOUT_GOT_CHECKOUT_SESSION = 'CHECKOUT_GOT_CHECKOUT_SESSION';

export const CHECKOUT_GOT_SUBSCRIPTIONS = 'CHECKOUT_GOT_SUBSCRIPTIONS';
export const CHECKOUT_CANCELED_SUBSCRIPTION = 'CHECKOUT_CANCELED_SUBSCRIPTION';

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

export function getCheckoutSession(id, productId) {
  return async dispatch => {
    dispatch(startLoading(id));
    try {
      const session = await makeCheckoutRequest('session', {
        productId
      });
      dispatch(gotCheckoutSession(session.data));
    } catch (e) {
      dispatch(checkoutError(id, e.response));
    }
    dispatch(endLoading(id));
  };
}

function gotSubscriptions(subscriptions) {
  return {
    type: CHECKOUT_GOT_SUBSCRIPTIONS,
    subscriptions
  };
}

export function getSubscriptions(id) {
  return async dispatch => {
    startLoading(id);
    try {
      const subsRequest = await makeCheckoutRequest('subscriptions');
      dispatch(gotSubscriptions(subsRequest.data));
    } catch (e) {
      dispatch(checkoutError(id, e.res));
    }
    endLoading(id);
  };
}

function canceledSubscription() {
  return {
    type: CHECKOUT_CANCELED_SUBSCRIPTION
  };
}

export function cancelSubscription(id) {
  return async dispatch => {
    startLoading(id);
    try {
      await makeCheckoutRequest('subscriptions', {}, 'delete');
      dispatch(canceledSubscription());
    } catch (e) {
      dispatch(checkoutError(id, e.res));
    }
    endLoading(id);
  };
}
