import { all, put, takeEvery } from 'redux-saga/effects';
import { endLoading, startLoading } from '../actions/loading';
import makeCanvasRequest from '../util/canvas/makeCanvasRequest';
import { canvasProxyError } from '../actions/error';
import {
  CANVAS_GET_OBSERVEES,
  CANVAS_GET_USER_COURSES,
  gotObservees,
  gotUserCourses
} from '../actions/canvas';
import getGradedUsersFromCourses from '../util/canvas/getGradedUsersFromCourses';

function* getUserCourses({ id, token, subdomain }) {
  yield put(startLoading(id));
  try {
    const userRes = yield makeCanvasRequest('courses', token, subdomain);
    yield put(
      gotUserCourses(userRes.data, getGradedUsersFromCourses(userRes.data))
    );
  } catch (e) {
    yield put(canvasProxyError(id, e.response));
  }
  yield put(endLoading(id));
}

function* getObservees({ id, token, subdomain, userId }) {
  yield put(startLoading(id));
  try {
    const observeesRequest = yield makeCanvasRequest(
      'users/profile/self/observees',
      token,
      subdomain,
      { user_id: userId }
    );
    yield put(gotObservees(observeesRequest.data));
  } catch (e) {
    yield put(canvasProxyError(id, e.res));
  }
  yield put(endLoading(id));
}

function* watcher() {
  yield takeEvery(CANVAS_GET_USER_COURSES, getUserCourses);
  yield takeEvery(CANVAS_GET_OBSERVEES, getObservees);
}

export default function* canvasRootSaga() {
  yield all([watcher()]);
}
