import makePlusRequest from '../util/plus/makePlusRequest';
import { plusError } from './error';
import { startLoading, endLoading } from './loading';

export const PLUS_GOT_SESSION_INFORMATION = 'PLUS_GOT_SESSION_INFORMATION';

export const PLUS_GOT_AVERAGE_GRADE_FOR_COURSE =
  'PLUS_GOT_AVERAGE_GRADE_FOR_COURSE';

export const PLUS_GOT_AVERAGE_SCORE_FOR_OUTCOME =
  'PLUS_GOT_AVERAGE_SCORE_FOR_OUTCOME';

function gotSessionInformation(sessionInformation) {
  return {
    type: PLUS_GOT_SESSION_INFORMATION,
    sessionInformation
  };
}

export function getSessionInformation(id) {
  return async dispatch => {
    startLoading(id);
    try {
      const sessionRequest = await makePlusRequest('session');
      dispatch(gotSessionInformation(sessionRequest.data));
    } catch (e) {
      dispatch(plusError(id, e.res));
    }
    endLoading(id);
  };
}

function gotAverageGradeForCourse(courseId, averageGrade) {
  return {
    type: PLUS_GOT_AVERAGE_GRADE_FOR_COURSE,
    courseId,
    averageGrade
  };
}

export function getAverageGradeForCourse(id, courseId) {
  return async dispatch => {
    startLoading(id);
    try {
      const avgGradeRequest = await makePlusRequest(`courses/${courseId}/avg`);
      dispatch(gotAverageGradeForCourse(courseId, avgGradeRequest.data));
    } catch (e) {
      dispatch(plusError(id, e.res));
    }
    endLoading(id);
  };
}

function gotAverageScoreForOutcome(outcomeId, avg) {
  return {
    type: PLUS_GOT_AVERAGE_SCORE_FOR_OUTCOME,
    outcomeId,
    avg
  };
}

export function getAverageScoreForOutcome(id, outcomeId) {
  return async dispatch => {
    startLoading(id);
    try {
      const avgOutcomeResponse = await makePlusRequest(
        `outcomes/${outcomeId}/avg`
      );
      dispatch(gotAverageScoreForOutcome(outcomeId, avgOutcomeResponse.data));
    } catch (e) {
      dispatch(plusError(id, e.res));
    }
    endLoading(id);
  };
}
