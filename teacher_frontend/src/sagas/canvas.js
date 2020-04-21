import { all, put, takeLeading, takeEvery } from "redux-saga/effects";
import makeApiRequest from "../util/api/makeApiRequest";
import {
  CANVAS_GET_COURSE_ENROLLMENTS,
  CANVAS_GET_COURSES,
  CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW,
  CANVAS_GET_USER_PROFILE,
  CANVAS_LOGOUT,
  getCourseEnrollmentsError,
  getCoursesError,
  getDistanceLearningGradesOverviewError,
  getUserProfileError,
  gotCourseEnrollments,
  gotCourses,
  gotDistanceLearningGradesOverview,
  gotUserProfile,
  loggedOut,
  logoutError
} from "../actions/canvas";
import makePlusRequest from "../util/plus/makePlusRequest";

function* getUserProfile() {
  try {
    const userProfileResponse = yield makeApiRequest(`users/self/profile`);
    yield put(gotUserProfile(userProfileResponse.data.profile));
  } catch (e) {
    yield put(
      getUserProfileError(e.response ? e.response.data : "can't connect")
    );
  }
}

function* getCourses() {
  try {
    const coursesResponse = yield makeApiRequest(`courses`, {
      include: ["distance_learning_pairs"]
    });
    yield put(
      gotCourses(
        coursesResponse.data.courses,
        coursesResponse.data.distance_learning_pairs
      )
    );
  } catch (e) {
    yield put(getCoursesError(e.response ? e.response.data : "can't connect"));
  }
}

function* getDistanceLearningGradesOverview({
  originalCourseId,
  distanceLearningCourseId
}) {
  try {
    const dlgOverviewResponse = yield makeApiRequest(
      `grades/distance_learning/overview`,
      {
        original_course_id: originalCourseId,
        distance_learning_course_id: distanceLearningCourseId
      }
    );
    yield put(
      gotDistanceLearningGradesOverview(
        originalCourseId,
        distanceLearningCourseId,
        dlgOverviewResponse.data.distance_learning_grades_overview
      )
    );
  } catch (e) {
    yield put(
      getDistanceLearningGradesOverviewError(
        e.response ? e.response.data : "can't connect"
      )
    );
  }
}

function* getCourseEnrollments({ courseId }) {
  try {
    const enrollmentsRequest = yield makeApiRequest(
      `courses/${courseId}/enrollments`,
      { type: ["StudentEnrollment"], state: ["active"] }
    );
    yield put(
      gotCourseEnrollments(courseId, enrollmentsRequest.data.enrollments)
    );
  } catch (e) {
    yield put(
      getCourseEnrollmentsError(e.response ? e.response.data : "can't connect")
    );
  }
}

function* logout() {
  try {
    yield makePlusRequest("session", {}, "delete");
    yield put(loggedOut());
  } catch (e) {
    put(logoutError(e.response ? e.response.data : "can't connect"));
  }
}

function* watcher() {
  yield takeLeading(CANVAS_GET_USER_PROFILE, getUserProfile);
  yield takeLeading(CANVAS_GET_COURSES, getCourses);
  yield takeEvery(
    CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW,
    getDistanceLearningGradesOverview
  );
  yield takeEvery(CANVAS_GET_COURSE_ENROLLMENTS, getCourseEnrollments);
  yield takeLeading(CANVAS_LOGOUT, logout);
}

export default function* canvasRootSaga() {
  yield all([watcher()]);
}
