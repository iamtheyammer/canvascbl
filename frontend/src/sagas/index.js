import { all } from 'redux-saga/effects';

import canvasRootSaga from './canvas';
import checkoutRootSaga from './checkout';

export default function* rootSaga() {
  yield all([canvasRootSaga(), checkoutRootSaga()]);
}
