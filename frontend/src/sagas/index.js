import { all } from 'redux-saga/effects';

import canvasRootSaga from './canvas';
import checkoutRootSaga from './checkout';
import plusRootSaga from './plus';
import oauth2RootSaga from './oauth2';

export default function* rootSaga() {
  yield all([
    canvasRootSaga(),
    checkoutRootSaga(),
    plusRootSaga(),
    oauth2RootSaga()
  ]);
}
