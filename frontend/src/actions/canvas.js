export const CANVAS_GET_INITIAL_DATA = 'CANVAS_GET_INITIAL_DATA';
export const CANVAS_GOT_INITIAL_DATA = 'CANVAS_GOT_INITIAL_DATA';

export const CANVAS_LOGOUT = 'CANVAS_LOGOUT';

export const CANVAS_DELETED_TOKEN = 'CANVAS_DELETED_TOKEN';

export const CANVAS_GOT_STORED_CREDENTIALS = 'CANVAS_GOT_STORED_CREDENTIALS';

export const CANVAS_GOT_TOKEN_ENTRY = 'CANVAS_GOT_TOKEN_ENTRY';
export const CANVAS_GOT_USER_OAUTH = 'CANVAS_GOT_USER_OAUTH';

export const CANVAS_SENT_TOKEN = 'CANVAS_SENT_TOKEN';

export const CANVAS_GOT_TOKEN = 'CANVAS_GOT_TOKEN';

// dep
export const CANVAS_GOT_NEW_TOKEN_FROM_REFRESH_TOKEN =
  'CANVAS_GOT_NEW_TOKEN_FROM_REFRESH_TOKEN';

export const CANVAS_GOT_USER_PROFILE = 'CANVAS_GOT_USER_PROFILE';

export const CANVAS_GET_USER_COURSES = 'CANVAS_GET_USER_COURSES';
export const CANVAS_GOT_USER_COURSES = 'CANVAS_GOT_USER_COURSES';

export const CANVAS_GOT_OUTCOME_RESULTS_FOR_COURSE =
  'CANVAS_GOT_OUTCOME_RESULTS_FOR_COURSE';

export const CANVAS_GOT_OUTCOME_ROLLUPS_FOR_COURSE =
  'CANVAS_GOT_OUTCOME_ROLLUPS_FOR_COURSE';

export const CANVAS_GOT_OUTCOME_ROLLUPS_AND_OUTCOMES_FOR_COURSE =
  'CANVAS_GOT_OUTCOME_ROLLUPS_AND_OUTCOMES_FOR_COURSE';

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
  grades
) {
  return {
    type: CANVAS_GOT_INITIAL_DATA,
    user,
    observees,
    courses,
    gradedUsers,
    outcomeResults,
    grades
  };
}

export function logout() {
  // even that these are deprecated, we'll keep it here to flush old stuff
  localStorage.clear();

  // super duper clear everything.
  window.location.reload();

  return {
    type: CANVAS_LOGOUT
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
