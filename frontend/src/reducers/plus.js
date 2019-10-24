import {
  PLUS_GOT_AVERAGE_GRADE_FOR_COURSE,
  PLUS_GOT_SESSION_INFORMATION
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
    default:
      return state;
  }
}
