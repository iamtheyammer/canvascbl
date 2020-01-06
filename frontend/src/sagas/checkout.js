import { all, put, takeLatest } from 'redux-saga/effects';
import { endLoading, startLoading } from '../actions/loading';
import makeCheckoutRequest from '../util/checkout/makeCheckoutRequest';
import {
  REDEEM_REDEEM_GIFT_CARDS,
  redeemGiftCardsSuccess
} from '../actions/components/redeem';
import { checkoutError } from '../actions/error';

function* redeemGiftCards({ id, giftCards }) {
  yield put(startLoading(id));
  try {
    const redeemResponse = yield makeCheckoutRequest('redeem', {}, 'post', {
      codes: giftCards
    });
    yield put(redeemGiftCardsSuccess(redeemResponse.data));
  } catch (e) {
    yield put(checkoutError(id, e.response && e.response.data));
  }
  yield put(endLoading(id));
}

function* watcher() {
  yield takeLatest(REDEEM_REDEEM_GIFT_CARDS, redeemGiftCards);
}

export default function* checkoutRootSaga() {
  yield all([watcher()]);
}
