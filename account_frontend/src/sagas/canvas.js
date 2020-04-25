import { all, put, takeLeading } from 'redux-saga/effects';
import makeApiRequest from '../util/api/makeApiRequest';
import {
  CANVAS_GET_USER_PROFILE,
  CANVAS_LOGOUT,
  getUserProfileError,
  gotUserProfile,
  loggedOut,
  logoutError
} from '../actions/canvas';
import makePlusRequest from '../util/plus/makePlusRequest';

function* getUserProfile() {
  try {
    const profileRequest = yield makeApiRequest('users/self/profile');
    yield put(gotUserProfile(profileRequest.data.profile));
  } catch (e) {
    yield put(
      getUserProfileError(e.response ? e.response.data : "can't connect")
    );
  }
}

function* logout() {
  try {
    yield makePlusRequest('session', {}, 'delete');
    yield put(loggedOut());
  } catch (e) {
    put(logoutError(e.response ? e.response.data : "can't connect"));
  }
}

function* watcher() {
  yield takeLeading(CANVAS_GET_USER_PROFILE, getUserProfile);
  yield takeLeading(CANVAS_LOGOUT, logout);
}

export default function* canvasRootSaga() {
  yield all([watcher()]);
}
