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

// Expand code list values based on the official documentation
// http://standard.open-contracting.org/latest/en/schema/codelists
function codeList(list, val) {
  let table = {};
  switch (list) {
    case 'tender.status':
      // http://standard.open-contracting.org/latest/en/schema/codelists/#tender-status
      table = {
        'planned':      'Planeada',
        'active':       'Activa',
        'cancelled':    'Cancelada',
        'unsuccessful': 'No exitosa',
        'complete':     'Cerrada'
      }
      break;
    case 'tender.procurementMethod':
      // http://standard.open-contracting.org/latest/en/schema/codelists/#method
      table = {
        'open':      'Licitación Pública',
        'selective': 'Invitación a 3 proveedores',
        'limited':   'Asignación Directa'
      }
      break;
    case 'award.status':
      // http://standard.open-contracting.org/latest/en/schema/codelists/#award-status
      table = {
        'pending':      'Pendiente',
        'active':       'Activa',
        'cancelled':    'Cancelada',
        'unsuccessful': 'No exitosa'
      }
      break;
    case 'contract.status':
      // http://standard.open-contracting.org/latest/en/schema/codelists/#contract-status
      table = {
        'pending':    'Pendiente',
        'active':     'Activo',
        'cancelled':  'Cancelado',
        'terminated': 'Finalizado'
      }
      break;
    default:
      return 'Invalid list';
  }
  return table[val];
}

export {
  getParameter,
  formatAmount,
  formatDate,
  codeList
}
