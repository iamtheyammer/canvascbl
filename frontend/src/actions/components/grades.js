export const GRADES_SWITCH_VIEW_TYPE = 'GRADES_SWITCH_VIEW_TYPE';

export function switchViewType(newTypeName) {
  return {
    type: GRADES_SWITCH_VIEW_TYPE,
    newTypeName
  };
}
