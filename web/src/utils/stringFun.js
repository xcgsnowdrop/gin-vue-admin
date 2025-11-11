/* eslint-disable */
export const toUpperCase = (str) => {
  if (str[0]) {
    return str.replace(str[0], str[0].toUpperCase())
  } else {
    return ''
  }
}

export const toLowerCase = (str) => {
  if (str[0]) {
    return str.replace(str[0], str[0].toLowerCase())
  } else {
    return ''
  }
}

// 驼峰转换下划线
export const toSQLLine = (str) => {
  if (str === 'ID') return 'ID'
  return str.replace(/([A-Z])/g, '_$1').toLowerCase()
}

// 下划线转换驼峰
export const toHump = (name) => {
  return name.replace(/\_(\w)/g, function (all, letter) {
    return letter.toUpperCase()
  })
}

/**
 * 将逗号分隔的字符串转换为整数数组
 * @param {string} str - 逗号分隔的字符串，例如 "1,2,3" 或 "1, 2, 3"
 * @returns {number[]} 整数数组，例如 [1, 2, 3]
 * @example
 * stringToIntArray("1,2,3") => [1, 2, 3]
 * stringToIntArray("1, 2, 3") => [1, 2, 3]
 * stringToIntArray("") => []
 * stringToIntArray("1,abc,3") => [1, 3]
 */
export const stringToIntArray = (str) => {
  if (!str || typeof str !== 'string' || str.trim() === '') {
    return []
  }
  return str
    .split(',')
    .map(item => item.trim())
    .filter(item => item !== '')
    .map(item => parseInt(item, 10))
    .filter(item => !isNaN(item))
}