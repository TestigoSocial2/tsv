import moment from 'moment';
import 'moment/locale/es';

// Helper method to retrieve a GET variable
function getParameter(name, url) {
  if (!url) {
    url = window.location.href;
  }
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
  results = regex.exec(url);
  if (!results) return false;
  if (!results[2]) return false;
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}

// Format a given value to a currency string
function formatAmount(val) {
  return '$' + val.toLocaleString(undefined, {
    useGrouping: true,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

// Pretty print date values
function formatDate( date, format ) {
  return moment(date).format( format );
}

export {
  getParameter,
  formatAmount,
  formatDate
}
