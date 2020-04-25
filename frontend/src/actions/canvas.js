export const CANVAS_GET_INITIAL_DATA = 'CANVAS_GET_INITIAL_DATA';
export const CANVAS_GOT_INITIAL_DATA = 'CANVAS_GOT_INITIAL_DATA';

export const CANVAS_GET_ASSIGNMENTS_FOR_COURSE =
  'CANVAS_GET_ASSIGNMENTS_FOR_COURSE';
export const CANVAS_GOT_ASSIGNMENTS_FOR_COURSE =
  'CANVAS_GOT_ASSIGNMENTS_FOR_COURSE';

export const CANVAS_GET_OUTCOME_ALIGNMENTS_FOR_COURSE =
  'CANVAS_GET_OUTCOME_ALIGNMENTS_FOR_COURSE';

export const CANVAS_GOT_OUTCOME_ALIGNMENTS_FOR_COURSE =
  'CANVAS_GOT_OUTCOME_ALIGNMENTS_FOR_COURSE';

export const CANVAS_GET_OBSERVEES = 'CANVAS_GET_OBSERVEES';
export const CANVAS_GOT_OBSERVEES = 'CANVAS_GOT_OBSERVEES';

export const CANVAS_GET_INDIVIDUAL_OUTCOME = 'CANVAS_GET_INDIVIDUAL_OUTCOME';
export const CANVAS_GOT_INDIVIDUAL_OUTCOME = 'CANVAS_GOT_INDIVIDUAL_OUTCOME';

export const CANVAS_CHANGE_ACTIVE_USER = 'CANVAS_CHANGE_ACTIVE_USER';

export const CANVAS_TOGGLE_COURSE_VISIBILITY =
  'CANVAS_TOGGLE_COURSE_VISIBILITY';
export const CANVAS_TOGGLE_COURSE_VISIBILITY_ERROR =
  'CANVAS_TOGGLE_COURSE_VISIBILITY_ERROR';
export const CANVAS_TOGGLED_COURSE_VISIBILITY =
  'CANVAS_TOGGLED_COURSE_VISIBILITY';

export function getInitialData(id) {
  return {
    type: CANVAS_GET_INITIAL_DATA,
    id
  };
}

export function gotInitialData(
  user,
  observees,
  courses,
  gradedUsers,
  outcomeResults,
  grades,
  gpa,
  distanceLearning
) {
  return {
    type: CANVAS_GOT_INITIAL_DATA,
    user,
    observees,
    courses,
    gradedUsers,
    outcomeResults,
    grades,
    gpa,
    distanceLearning
  };
}

export function gotAssignmentsForCourse(assignments, courseId) {
  return {
    type: CANVAS_GOT_ASSIGNMENTS_FOR_COURSE,
    assignments,
    courseId
  };
}

export function getAssignmentsForCourse(id, courseId) {
  return {
    type: CANVAS_GET_ASSIGNMENTS_FOR_COURSE,
    id,
    courseId
  };
}

export function gotOutcomeAlignmentsForCourse(courseId, alignments) {
  return {
    type: CANVAS_GOT_OUTCOME_ALIGNMENTS_FOR_COURSE,
    alignments,
    courseId
  };
}

export function getOutcomeAlignmentsForCourse(id, courseId, studentId) {
  return {
    type: CANVAS_GET_OUTCOME_ALIGNMENTS_FOR_COURSE,
    id,
    courseId,
    studentId
  };
}

export function getIndividualOutcome(id, outcomeId) {
  return {
    type: CANVAS_GET_INDIVIDUAL_OUTCOME,
    id,
    outcomeId
  };
}

export function gotIndividualOutcome(outcome) {
  return {
    type: CANVAS_GOT_INDIVIDUAL_OUTCOME,
    outcome
  };
}

export function changeActiveUser(id) {
  return {
    type: CANVAS_CHANGE_ACTIVE_USER,
    id
  };
}

export function toggleCourseVisibility(id, courseId, toggle) {
  return {
    type: CANVAS_TOGGLE_COURSE_VISIBILITY,
    id,
    courseId,
    toggle
  };
}

export function toggleCourseVisibilityError(e, courseId) {
  return {
    type: CANVAS_TOGGLE_COURSE_VISIBILITY_ERROR,
    e,
    courseId
  };
}

export function toggledCourseVisibility(courseId, newStatus) {
  return {
    type: CANVAS_TOGGLED_COURSE_VISIBILITY,
    courseId,
    newStatus
  };
}
