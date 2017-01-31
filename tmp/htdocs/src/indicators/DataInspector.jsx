import React from 'react';
import Chart from 'chart.js'

class DataInspector extends React.Component {
  constructor(props) {
    super(props);
    this.onChange = this.onChange.bind(this);
    this.applyFilters = this.applyFilters.bind(this);
    this.state = {
      filters: props.defaultFilters
    };
  }

  onChange(e) {
    let name = e.target.name;
    this.setState({
      filters: Object.assign(this.state.filters, {[name]: e.target.value})
    });
  }

  componentDidMount() {
    // Setup slider
    let amountSlider = $("#amountSlider");
    amountSlider.slider({
      step: 50000,
      min: 0,
      max: 500000000,
      value: [20000000,80000000],
      formatter: function( value ) {
        if( Array.isArray( value ) ) {
          var lbl = '$' + value[0].toLocaleString() + ' a ' + '$' + value[1].toLocaleString();
          return lbl;
        }
        return '';
      }
    }).on('change', (e) => (
      this.setState({ filters: Object.assign(this.state.filters, {amount:e.value.newValue}) })
    ));
  }

  applyFilters() {
    this.props.onSubmit(this.state.filters);
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
                <input type="radio" name="bucket" value="gacm" onChange={this.onChange} checked={this.state.filters.bucket == 'gacm'} />Grupo Aeroportuario
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="bucket" value="cdmx" onChange={this.onChange} checked={this.state.filters.bucket == 'cdmx'} />Ciudad de México
              </label>
            </div>
            <hr />
            <h4>Etapa del Procedimiento</h4>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="planning" onChange={this.onChange} checked={this.state.filters.state == 'planning'} />Planeación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="tender" onChange={this.onChange} checked={this.state.filters.state == 'tender'} />Licitación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="award" onChange={this.onChange} checked={this.state.filters.state == 'award'} />Adjudicación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="contract" onChange={this.onChange} checked={this.state.filters.state == 'contract'} />Contratación
              </label>
            </div>
            <div className="radio">
              <label>
                <input type="radio" name="state" value="implementation" onChange={this.onChange} checked={this.state.filters.state == 'implementation'} />Implementación
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
          <h2>Resultados de la Búsqueda</h2>
          <div className="chart">
            <canvas></canvas>
          </div>
        </div>
      </div>
    );
  }
}

export default DataInspector;
