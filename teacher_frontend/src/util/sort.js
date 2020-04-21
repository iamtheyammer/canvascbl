import moment from "moment";

const sort = (a, b) => {
  const al = a.toLocaleLowerCase();
  const bl = b.toLocaleLowerCase();
  if (al < bl) return -1;
  if (al > bl) return 1;
  return 0;
};

export const desc = sort;

export const asc = (a, b) => sort(b, a);

const date = (a, b) => {
  const ma = moment(a);
  const mb = moment(b);
  if (ma.isBefore(mb)) return -1;
  if (ma.isAfter(mb)) return 1;
  return 0;
};

export const dateDesc = date;

export const dateAsc = (a, b) => date(b, a);
