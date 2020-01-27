import { put, takeEvery, all } from 'redux-saga/effects';
import { endLoading, startLoading } from '../actions/loading';
import makePlusRequest from '../util/plus/makePlusRequest';
import { CANVAS_LOGOUT, loggedOut } from '../actions/canvas';
import { plusError } from '../actions/error';

function* logout({ id }) {
  yield put(startLoading(id));
  try {
    yield makePlusRequest('session', {}, 'delete');
    yield put(loggedOut(id));
  } catch (e) {
    yield put(plusError(id, e.response));
  }
  yield put(endLoading(id));
}

function* watcher() {
  yield takeEvery(CANVAS_LOGOUT, logout);
}

export default function* plusRootSaga() {
  yield all([watcher()]);
}
