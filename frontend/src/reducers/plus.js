import {
  PLUS_GOT_AVERAGE_GRADE_FOR_COURSE,
  PLUS_GOT_SESSION_INFORMATION,
  PLUS_GOT_AVERAGE_SCORE_FOR_OUTCOME,
  PLUS_GOT_PREVIOUS_GRADES
} from '../actions/plus';

export default function plus(state = [], action) {
  switch (action.type) {
    case PLUS_GOT_SESSION_INFORMATION:
      return {
        ...state,
        session: action.sessionInformation
      };
    case PLUS_GOT_AVERAGE_GRADE_FOR_COURSE:
      return {
        ...state,
        averages: {
          ...state.averages,
          [action.courseId]: action.averageGrade
        }
      };
    case PLUS_GOT_AVERAGE_SCORE_FOR_OUTCOME:
      return {
        ...state,
        outcomeAverages: {
          ...state.outcomeAverages,
          [action.outcomeId]: action.avg
        }
      };
    case PLUS_GOT_PREVIOUS_GRADES:
      return {
        ...state,
        previousGrades: action.previousGrades
      };
    default:
      return state;
  }
}
