export const FILTERS_CLEAR = 'FILTERS_CLEAR';
export const FILTERS_UPDATE_NAME = 'FILTERS_UPDATE_NAME';
export const FILTERS_UPDATE_NAME_TYPE = 'FILTERS_UPDATE_NAME_TYPE';

export function clearFilters() {
  return {
    type: FILTERS_CLEAR,
  };
}

export function filterName(newName) {
  return {
    type: FILTERS_UPDATE_NAME,
    newName,
  };
}

export function filterNameType(newType) {
  return {
    type: FILTERS_UPDATE_NAME_TYPE,
    newType,
  };
}
