import {
  CANVAS_GOT_USER_OAUTH,
  CANVAS_GOT_TOKEN_ENTRY,
  CANVAS_GOT_NEW_TOKEN_FROM_REFRESH_TOKEN,
  CANVAS_GOT_USER_PROFILE,
  CANVAS_GOT_USER_COURSES,
  CANVAS_GOT_OUTCOME_ROLLUPS_FOR_COURSE,
  CANVAS_GOT_OUTCOME_RESULTS_FOR_COURSE,
  CANVAS_GOT_OUTCOME_ROLLUPS_AND_OUTCOMES_FOR_COURSE,
  CANVAS_GOT_ASSIGNMENTS_FOR_COURSE,
  CANVAS_GOT_STORED_CREDENTIALS,
  CANVAS_GOT_TOKEN,
  CANVAS_SENT_TOKEN,
  CANVAS_DELETED_TOKEN,
  CANVAS_GOT_OUTCOME_ALIGNMENTS_FOR_COURSE,
  CANVAS_GOT_OBSERVEES,
  CANVAS_CHANGE_ACTIVE_USER,
  CANVAS_GOT_INITIAL_DATA,
  CANVAS_GOT_INDIVIDUAL_OUTCOME
} from '../actions/canvas';

export default function canvas(state = {}, action) {
  switch (action.type) {
    case CANVAS_DELETED_TOKEN:
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
    case CANVAS_GOT_STORED_CREDENTIALS:
      return {
        ...state,
        ...{
          token: action.token,
          refreshToken: action.refreshToken,
          subdomain: action.subdomain
        }
      };
    case CANVAS_GOT_TOKEN_ENTRY:
      return {
        ...state,
        ...{
          token: action.token,
          subdomain: action.subdomain
        }
      };
    case CANVAS_GOT_USER_OAUTH:
      return {
        ...state,
        token: action.token,
        refreshToken: action.refreshToken,
        subdomain: action.subdomain
      };
    case CANVAS_GOT_NEW_TOKEN_FROM_REFRESH_TOKEN:
      return {
        ...state,
        token: action.newToken
      };
    case CANVAS_SENT_TOKEN:
      return {
        ...state,
        successfullySentToken: true
      };
    case CANVAS_GOT_TOKEN:
      return {
        ...state,
        token: action.token,
        subdomain: action.subdomain
      };
    case CANVAS_GOT_USER_PROFILE:
      return {
        ...state,
        user: action.user,
        users: {
          ...state.users,
          [action.user.id]: action.user
        }
      };
    case CANVAS_GOT_USER_COURSES:
      return {
        ...state,
        courses: action.courses,
        gradedUsers: action.gradedUsers,
        // by default, the active user is the first graded user
        activeUserId: action.gradedUsers[0]
      };
    case CANVAS_GOT_OUTCOME_ROLLUPS_FOR_COURSE:
      const courseId = action.courseId;
      return {
        ...state,
        ...{
          outcomeRollups: {
            ...state.outcomeRollups,
            ...{
              [courseId]: action.results
            }
          }
        }
      };
    case CANVAS_GOT_OUTCOME_RESULTS_FOR_COURSE:
      return {
        ...state,
        ...{
          outcomeResults: {
            ...state.outcomeResults,
            [action.courseId]: action.results
          },
          outcomes: {
            ...state.outcomes
          }
        }
      };
    case CANVAS_GOT_OUTCOME_ROLLUPS_AND_OUTCOMES_FOR_COURSE:
      return {
        ...state,
        ...{
          outcomeRollups: {
            ...state.outcomeRollups,
            [action.courseId]: action.results
          },
          outcomes: {
            ...state.outcomes,
            [action.courseId]: action.outcomes
          }
        }
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
