export const COURSESETTINGS_SET_TOGGLE_COURSE_VISIBILITY_ID =
  'COURSESETTINGS_SET_TOGGLE_COURSE_VISIBILITY_ID';

export function setToggleCourseVisibilityId(id, courseId) {
  return {
    type: COURSESETTINGS_SET_TOGGLE_COURSE_VISIBILITY_ID,
    id,
    courseId
  };
}
