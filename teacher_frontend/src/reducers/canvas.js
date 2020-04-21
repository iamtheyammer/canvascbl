import {
  CANVAS_GET_COURSE_ENROLLMENTS,
  CANVAS_GET_COURSE_ENROLLMENTS_ERROR,
  CANVAS_GET_COURSES,
  CANVAS_GET_COURSES_ERROR,
  CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW,
  CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW_ERROR,
  CANVAS_GOT_COURSE_ENROLLMENTS,
  CANVAS_GOT_COURSES,
  CANVAS_GOT_DISTANCE_LEARNING_GRADES_OVERVIEW
} from "../actions/canvas";

export default function canvas(state = {}, action) {
  switch (action.type) {
    case CANVAS_GET_COURSES:
      return {
        ...state,
        loadingCourses: true
      };
    case CANVAS_GET_COURSES_ERROR:
      return {
        ...state,
        getCoursesError: action.e,
        loadingCourses: false
      };
    case CANVAS_GOT_COURSES:
      return {
        ...state,
        courses: action.courses,
        distanceLearningPairs: action.distanceLearningPairs,
        loadingCourses: false
      };
    case CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW:
      return {
        ...state,
        loadingDistanceLearningGradesOverview: true
      };
    case CANVAS_GET_DISTANCE_LEARNING_GRADES_OVERVIEW_ERROR:
      return {
        ...state,
        loadingDistanceLearningGradesOverview: false,
        getDistanceLearningGradesOverviewError: action.e
      };
    case CANVAS_GOT_DISTANCE_LEARNING_GRADES_OVERVIEW:
      return {
        ...state,
        loadingDistanceLearningGradesOverview: false,
        distanceLearningGradesOverviews: {
          ...state.distanceLearningGradesOverviews,
          [`${action.originalCourseId}_${action.distanceLearningCourseId}`]: action.overview
        }
      };
    case CANVAS_GET_COURSE_ENROLLMENTS:
      return {
        ...state,
        courseEnrollmentsAreLoading: true
      };
    case CANVAS_GET_COURSE_ENROLLMENTS_ERROR:
      return {
        ...state,
        getCourseEnrollmentsError: action.e,
        courseEnrollmentsAreLoading: false
      };
    case CANVAS_GOT_COURSE_ENROLLMENTS:
      return {
        ...state,
        enrollments: {
          ...state.enrollments,
          [action.courseId]: action.enrollments
        },
        courseEnrollmentsAreLoading: false
      };
    default:
      return state;
  }
}
