import {
  CHECKOUT_CANCELED_SUBSCRIPTION,
  CHECKOUT_GOT_CHECKOUT_SESSION,
  CHECKOUT_GOT_PRODUCTS,
  CHECKOUT_GOT_SUBSCRIPTIONS
} from '../actions/checkout';

export default function checkout(state = {}, action) {
  switch (action.type) {
    case CHECKOUT_GOT_PRODUCTS:
      return {
        ...state,
        products: action.products
      };
    case CHECKOUT_GOT_CHECKOUT_SESSION:
      return {
        ...state,
        ...{
          session: action.session
        }
      };
    case CHECKOUT_GOT_SUBSCRIPTIONS:
      return {
        ...state,
        subscriptions: action.subscriptions
      };
    case CHECKOUT_CANCELED_SUBSCRIPTION:
      return state;
    default:
      return state;
  }
}
