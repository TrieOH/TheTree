export function timeAgo(isoDate: string): string {
  const diff = Date.now() - new Date(isoDate).getTime();
  const sec = Math.floor(diff / 1000);
  const min = Math.floor(sec / 60);
  const hour = Math.floor(min / 60);
  const day = Math.floor(hour / 24);
  const week = Math.floor(day / 7);

  if (sec < 60) return "agora mesmo";
  if (min < 60) return `${min} minutos atrás`;
  if (hour < 24) return `${hour} horas atrás`;
  if (day < 7) return `${day} dias atrás`;
  return `${week} semanas atrás`;
}
