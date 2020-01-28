import {
  CANVAS_LOGOUT,
  CANVAS_GOT_ASSIGNMENTS_FOR_COURSE,
  CANVAS_GOT_OUTCOME_ALIGNMENTS_FOR_COURSE,
  CANVAS_GOT_OBSERVEES,
  CANVAS_CHANGE_ACTIVE_USER,
  CANVAS_GOT_INITIAL_DATA,
  CANVAS_GOT_INDIVIDUAL_OUTCOME
} from '../actions/canvas';

export default function canvas(state = {}, action) {
  switch (action.type) {
    case CANVAS_LOGOUT:
      return {};
    case CANVAS_GOT_INITIAL_DATA:
      return {
        ...state,
        user: action.user,
        users: {
          ...state.users,
          [action.user.id]: action.user,
          ...action.observees.reduce(
            (acc, val) => ({ ...acc, [val.id]: val }),
            {}
          )
        },
        observees: action.observees,
        courses: action.courses,
        gradedUsers: action.gradedUsers,
        // by default, the active user is the first graded user
        activeUserId: action.gradedUsers[0],
        outcomeResults: action.outcomeResults,
        grades: action.grades
      };
    case CANVAS_GOT_ASSIGNMENTS_FOR_COURSE:
      return {
        ...state,
        ...{
          assignments: {
            ...state.assignments,
            [action.courseId]: action.assignments
          }
        }
      };
    case CANVAS_GOT_OUTCOME_ALIGNMENTS_FOR_COURSE:
      return {
        ...state,
        outcomeAlignments: {
          ...state.outcomeAlignments,
          [action.courseId]: action.alignments
        }
      };
    case CANVAS_GOT_INDIVIDUAL_OUTCOME:
      return {
        ...state,
        outcomes: state.outcomes
          ? [...state.outcomes, action.outcome]
          : [action.outcome]
      };
    case CANVAS_GOT_OBSERVEES:
      return {
        ...state,
        observees: action.observees,
        users: {
          ...state.users,
          ...action.observees.reduce(
            (acc, val) => ({ ...acc, [val.id]: val }),
            {}
          )
        }
      };
    case CANVAS_CHANGE_ACTIVE_USER:
      return {
        ...state,
        activeUserId: action.id
      };
    default:
      return state;
  }
}
