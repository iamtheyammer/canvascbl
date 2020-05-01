import { all } from 'redux-saga/effects';

function* watcher() {}

export default function* plusRootSaga() {
  yield all([watcher()]);
}
