/**
 * 时间戳工具函数
 * 统一处理时间戳（秒级）与 Date 对象之间的转换
 */

/**
 * 将时间戳（秒级）转换为 Date 对象
 * @param {number|string|null|undefined} timestamp - 时间戳（秒级）
 * @returns {Date|null} Date 对象，如果输入无效则返回 null
 */
export function timestampToDate(timestamp) {
  if (!timestamp && timestamp !== 0) {
    return null
  }
  
  // 如果是 Date 对象，直接返回
  if (timestamp instanceof Date) {
    return timestamp
  }
  
  // 转换为数字
  const ts = typeof timestamp === 'string' ? Number(timestamp) : timestamp
  
  // 验证是否为有效数字
  if (isNaN(ts)) {
    return null
  }
  
  // 秒级时间戳转换为毫秒（乘以 1000）
  return new Date(ts * 1000)
}

/**
 * 将 Date 对象转换为时间戳（秒级）
 * @param {Date|string|number|null|undefined} date - Date 对象、时间字符串或时间戳
 * @returns {number|null} 时间戳（秒级），如果输入无效则返回 null
 */
export function dateToTimestamp(date) {
  if (!date) {
    return null
  }
  
  // 如果已经是时间戳（秒级），直接返回
  if (typeof date === 'number') {
    // 判断是秒级还是毫秒级时间戳
    // 毫秒级时间戳通常大于 10^10
    if (date > 10000000000) {
      return Math.floor(date / 1000) // 转换为秒级
    }
    return date
  }
  
  // 如果是 Date 对象，转换为时间戳
  if (date instanceof Date) {
    return Math.floor(date.getTime() / 1000)
  }
  
  // 如果是字符串，尝试转换为 Date
  if (typeof date === 'string') {
    const d = new Date(date)
    if (!isNaN(d.getTime())) {
      return Math.floor(d.getTime() / 1000)
    }
  }
  
  return null
}

/**
 * 格式化时间戳为字符串
 * @param {number|string|null|undefined|Date} timestamp - 时间戳（秒级）或 Date 对象
 * @param {string} format - 格式化模式，默认为本地化字符串
 * @returns {string} 格式化后的时间字符串，如果无效则返回 '-'
 */
export function formatTimestamp(timestamp, format = 'localeString') {
  const date = timestampToDate(timestamp)
  if (!date) {
    return '-'
  }
  
  switch (format) {
    case 'localeString':
      return date.toLocaleString()
    case 'localeDateString':
      return date.toLocaleDateString()
    case 'localeTimeString':
      return date.toLocaleTimeString()
    default:
      return date.toLocaleString()
  }
}

/**
 * 批量转换时间戳字段为 Date 对象
 * @param {Object} data - 数据对象
 * @param {string[]} fields - 需要转换的字段名数组
 * @returns {Object} 转换后的数据对象
 */
export function convertTimestampsToDates(data, fields = []) {
  if (!data || !Array.isArray(fields) || fields.length === 0) {
    return data
  }
  
  const result = { ...data }
  fields.forEach(field => {
    if (result[field] !== undefined) {
      result[field] = timestampToDate(result[field])
    }
  })
  
  return result
}

/**
 * 批量转换 Date 对象字段为时间戳（秒级）
 * @param {Object} data - 数据对象
 * @param {string[]} fields - 需要转换的字段名数组
 * @returns {Object} 转换后的数据对象
 */
export function convertDatesToTimestamps(data, fields = []) {
  if (!data || !Array.isArray(fields) || fields.length === 0) {
    return data
  }
  
  const result = { ...data }
  fields.forEach(field => {
    if (result[field] !== undefined && result[field] !== null) {
      result[field] = dateToTimestamp(result[field])
    }
  })
  
  return result
}

/**
 * 批量添加格式化时间字段
 * @param {Object} data - 数据对象
 * @param {Array<{source: string, target: string, format?: string}>} mappings - 字段映射配置
 * @returns {Object} 添加了格式化时间字段的数据对象
 */
export function addFormattedTimeFields(data, mappings = []) {
  if (!data || !Array.isArray(mappings) || mappings.length === 0) {
    return data
  }
  
  const result = { ...data }
  mappings.forEach(({ source, target, format = 'localeString' }) => {
    if (result[source] !== undefined) {
      result[target] = formatTimestamp(result[source], format)
    }
  })
  
  return result
}

