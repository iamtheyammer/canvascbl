import { all } from 'redux-saga/effects';

import canvasRootSaga from './canvas';
import checkoutRootSaga from './checkout';
import plusRootSaga from './plus';

export default function* rootSaga() {
  yield all([canvasRootSaga(), checkoutRootSaga(), plusRootSaga()]);
}
