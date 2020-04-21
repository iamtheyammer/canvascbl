export const CANVAS_GET_USER_PROFILE = 'CANVAS_GET_USER_PROFILE';
export const CANVAS_GET_USER_PROFILE_ERROR = 'CANVAS_GET_USER_PROFILE_ERROR';
export const CANVAS_GOT_USER_PROFILE = 'CANVAS_GOT_USER_PROFILE';

export const CANVAS_GET_COURSES = 'CANVAS_GET_COURSES';
export const CANVAS_GET_COURSES_ERROR = 'CANVAS_GET_COURSES_ERROR';
export const CANVAS_GOT_COURSES = 'CANVAS_GOT_COURSES';

export const CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW =
  'CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW';
export const CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW_ERROR =
  'CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW_ERROR';
export const CANVAS_GOT_DISTANCE_LEARNING_GRADES_OVERVIEW =
  'CANVAS_GOT_DISTANCE_LEARNING_GRADES_OVERVIEW';

export const CANVAS_GET_COURSE_ENROLLMENTS = 'CANVAS_GET_COURSE_ENROLLMENTS';
export const CANVAS_GET_COURSE_ENROLLMENTS_ERROR =
  'CANVAS_GET_COURSE_ENROLLMENTS_ERROR';
export const CANVAS_GOT_COURSE_ENROLLMENTS = 'CANVAS_GOT_COURSE_ENROLLMENTS';

export const CANVAS_LOGOUT = 'CANVAS_LOGOUT';
export const CANVAS_LOGOUT_ERROR = 'CANVAS_LOGOUT_ERROR';
export const CANVAS_LOGGED_OUT = 'CANVAS_LOGGED_OUT';

export function getUserProfile() {
  return {
    type: CANVAS_GET_USER_PROFILE,
  };
}

export function getUserProfileError(e) {
  return {
    type: CANVAS_GET_USER_PROFILE_ERROR,
    e,
  };
}

export function gotUserProfile(profile) {
  return {
    type: CANVAS_GOT_USER_PROFILE,
    profile,
  };
}

export function getCourses() {
  return {
    type: CANVAS_GET_COURSES,
  };
}

export function getCoursesError(e) {
  return {
    type: CANVAS_GET_COURSES_ERROR,
    e,
  };
}

export function gotCourses(courses, distanceLearningPairs) {
  return {
    type: CANVAS_GOT_COURSES,
    courses,
    distanceLearningPairs,
  };
}

export function getDistanceLearningGradesOverview(
  originalCourseId,
  distanceLearningCourseId
) {
  return {
    type: CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW,
    originalCourseId,
    distanceLearningCourseId,
  };
}

export function getDistanceLearningGradesOverviewError(e) {
  return {
    type: CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW_ERROR,
    e,
  };
}

export function gotDistanceLearningGradesOverview(
  originalCourseId,
  distanceLearningCourseId,
  overview
) {
  return {
    type: CANVAS_GOT_DISTANCE_LEARNING_GRADES_OVERVIEW,
    originalCourseId,
    distanceLearningCourseId,
    overview,
  };
}

export function getCourseEnrollments(courseId) {
  return {
    type: CANVAS_GET_COURSE_ENROLLMENTS,
    courseId,
  };
}

export function getCourseEnrollmentsError(e) {
  return {
    type: CANVAS_GET_COURSE_ENROLLMENTS_ERROR,
    e,
  };
}

export function gotCourseEnrollments(courseId, enrollments) {
  return {
    type: CANVAS_GOT_COURSE_ENROLLMENTS,
    courseId,
    enrollments,
  };
}

export function logout() {
  return {
    type: CANVAS_LOGOUT,
  };
}

export function logoutError(e) {
  return {
    type: CANVAS_LOGOUT_ERROR,
    e,
  };
}

export function loggedOut() {
  return {
    type: CANVAS_LOGGED_OUT,
  };
}
