import { all, put, takeEvery, takeLeading } from 'redux-saga/effects';
import { endLoading, startLoading } from '../actions/loading';
import { canvasProxyError } from '../actions/error';
import {
  CANVAS_GET_ASSIGNMENTS_FOR_COURSE,
  CANVAS_GET_INDIVIDUAL_OUTCOME,
  CANVAS_GET_INITIAL_DATA,
  CANVAS_GET_OUTCOME_ALIGNMENTS_FOR_COURSE,
  CANVAS_TOGGLE_COURSE_VISIBILITY,
  gotAssignmentsForCourse,
  gotIndividualOutcome,
  gotInitialData,
  gotOutcomeAlignmentsForCourse,
  toggleCourseVisibilityError,
  toggledCourseVisibility
} from '../actions/canvas';
import makeApiRequest from '../util/api/makeApiRequest';
import { gotSessionInformation } from '../actions/plus';

function* getInitialData({ id }) {
  yield put(startLoading(id));
  try {
    const gradesRequest = yield makeApiRequest('grades', {
      include: [
        'session',
        'user_profile',
        'observees',
        'courses',
        'outcome_results',
        'detailed_grades',
        'gpa'
      ]
    });
    const {
      session,
      user_profile,
      observees,
      courses,
      outcome_results,
      detailed_grades,
      gpa
    } = gradesRequest.data;
    yield put(gotSessionInformation(session));
    yield put(
      gotInitialData(
        user_profile,
        observees,
        courses,
        Object.keys(detailed_grades).map(uID => parseInt(uID)),
        outcome_results,
        detailed_grades,
        gpa
      )
    );
  } catch (e) {
    yield put(canvasProxyError(id, e.response));
  }
  yield put(endLoading(id));
}

function* getIndividualOutcome({ id, outcomeId }) {
  yield put(startLoading(id));
  try {
    const outcomeResponse = yield makeApiRequest(`outcomes/${outcomeId}`);
    yield put(gotIndividualOutcome(outcomeResponse.data));
  } catch (e) {
    yield put(canvasProxyError(id, e.response));
  }
  yield put(endLoading(id));
}

function* getOutcomeAlignmentsForCourse({ id, courseId, studentId }) {
  yield put(startLoading(id));
  try {
    const alignmentsResponse = yield makeApiRequest(
      `courses/${courseId}/outcome_alignments`,
      { student_id: studentId }
    );
    yield put(gotOutcomeAlignmentsForCourse(courseId, alignmentsResponse.data));
  } catch (e) {
    yield put(canvasProxyError(id, e.response));
  }
  yield put(endLoading(id));
}

function* getAssignmentsForCourse({ id, courseId }) {
  yield put(startLoading(id));
  try {
    const assignmentsResponse = yield makeApiRequest(
      `courses/${courseId}/assignments`
    );
    yield put(gotAssignmentsForCourse(assignmentsResponse.data, courseId));
  } catch (e) {
    yield put(canvasProxyError(id, e.response));
  }
  yield put(endLoading(id));
}

function* toggleCourseVisibility({ id, courseId, toggle }) {
  yield put(startLoading(id));
  try {
    const toggleResponse = yield makeApiRequest(
      `courses/${courseId}/hide`,
      {},
      toggle ? 'PUT' : 'DELETE'
    );
    yield put(toggledCourseVisibility(courseId, toggleResponse.data));
  } catch (e) {
    yield put(
      toggleCourseVisibilityError(
        e.response ? e.response.data : "can't connect",
        courseId
      )
    );
  }
  yield put(endLoading(id));
}

function* watcher() {
  yield takeLeading(CANVAS_GET_INITIAL_DATA, getInitialData);
  yield takeEvery(CANVAS_GET_INDIVIDUAL_OUTCOME, getIndividualOutcome);
  yield takeEvery(
    CANVAS_GET_OUTCOME_ALIGNMENTS_FOR_COURSE,
    getOutcomeAlignmentsForCourse
  );
  yield takeEvery(CANVAS_GET_ASSIGNMENTS_FOR_COURSE, getAssignmentsForCourse);
  yield takeLeading(CANVAS_TOGGLE_COURSE_VISIBILITY, toggleCourseVisibility);
}

export default function* canvasRootSaga() {
  yield all([watcher()]);
}
