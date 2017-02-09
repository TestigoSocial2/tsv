import React from 'react';
import ChartWidget from './ChartWidget.jsx';

class Details extends React.Component {
  constructor(props) {
    super(props);
    this.close = this.close.bind(this);
    this.applyFilters = this.applyFilters.bind(this);
  }

  close(e) {
    e.preventDefault();
    this.props.onClose();
  }

  componentDidMount() {
    // Set form reference and populate filter values
    this.form = $('form#filterForm');
    for( let k in this.props.filters ) {
      let selector = 'input[type="radio"][name="'+k+'"][value="'+this.props.filters[k]+'"]';
      this.form.find(selector).prop("checked",true);
    }

    // Setup slider
    let amountSlider = $("#amountSlider");
    amountSlider.slider({
      step: 50000,
      min: 0,
      max: 500000000,
      value: this.props.filters.amount,
      formatter: function( value ) {
        if( Array.isArray( value ) ) {
          var lbl = '$' + value[0].toLocaleString() + ' a ' + '$' + value[1].toLocaleString();
          return lbl;
        }
        return '';
      }
    });
  }

  applyFilters() {
    let data = {}
    this.form.serializeArray().forEach(function(el) {
      data[el.name] = el.value;
    });
    data.amount = data.amount.split(',');
    data.amount[0] = parseInt(data.amount[0]);
    data.amount[1] = parseInt(data.amount[1]);
    this.props.onSubmit(data);
  }

  render() {
    return (
      <div className="row chart-main">
        <div className="col-md-4 sidebar">
          <form id="filterForm">
            <h4>Filtros</h4>
            <p>Selecciona los criterios que te sean más relevantes para ajustar el indicador.</p>
            <hr />
            <h4>Dependencias</h4>
            <div className="radio">
              <label>
                <input type="radio" name="bucket" value="gacm" />Grupo Aeroportuario
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="bucket" value="cdmx" />Ciudad de México
              </label>
            </div>
            <hr />
            <h4>Etapa del Procedimiento</h4>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="planning" />Planeación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="tender" />Licitación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="award" />Adjudicación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="contract" />Contratación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="implementation" />Implementación
              </label>
            </div>
            <hr />
            <h4>Monto</h4>
            <b>$0</b><input id="amountSlider" name="amount" type="text" /><b>$500,000,000</b>
            <hr />
            <span className="btn-black active" onClick={this.applyFilters}>Aplicar Filtros</span>
          </form>
        </div>
        <div className="col-md-8 content">
          <h2>
            <button onClick={this.close} type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>  Detalles del Indicador
          </h2>

          <ChartWidget
            id="indicatorDetails"
            data={this.props.data}
            reducer={this.props.reducer}
            width="680"
            height="440" />
        </div>
      </div>
    );
  }
}

export default Details;
