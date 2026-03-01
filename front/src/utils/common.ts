export function numberToSizeStr(size: number, template?: number): string {
  if (size < 0) {
    return "-";
  }

  const KB = 1024;
  const MB = KB * 1024;
  const GB = MB * 1024;
  if (template!==undefined) {
    if (template! >= GB) {
      return (size / GB).toFixed(2) + " GB";
    } else if (template! >= MB) {
      return (size / MB).toFixed(2) + " MB";
    } else if (template! >= KB) {
      return (size / KB).toFixed(2) + " KB";
    } else {
      return size + " B"
    }
  }
  if (size >= GB) {
    return (size / GB).toFixed(2) + " GB";
  } else if (size >= MB) {
    return (size / MB).toFixed(2) + " MB";
  } else if (size >= KB) {
    return (size / KB).toFixed(2) + " KB";
  } else {
    return size + " B";
  }
}
