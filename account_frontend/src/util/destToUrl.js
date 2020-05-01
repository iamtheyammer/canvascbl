import env from './env';

export default function (dest) {
  switch (dest) {
    case 'canvascbl':
      return env.canvascblUrl;
    case 'teacher':
      return env.teacherUrl;
    default:
      return env.canvascblUrl;
  }
}

export function validateDest(dest) {
  switch (dest) {
    case 'canvascbl':
      return true;
    case 'teacher':
      return true;
    default:
      return false;
  }
}
