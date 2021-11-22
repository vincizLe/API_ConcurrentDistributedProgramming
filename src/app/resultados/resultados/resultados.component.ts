import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ServicesService } from '../services/services.service';

@Component({
  selector: 'app-resultados',
  templateUrl: './resultados.component.html',
  styleUrls: ['./resultados.component.scss'],
})
export class ResultadosComponent implements OnInit {
  public form!: FormGroup;
  public selectedFabricante: any = 0;
  public porcentaje: any = 0;
  public porcentajeTotal: any = 1;
  public isResultado = false;
  //grafico
  public chartDatasets:Array<any> = [{ data: [50, 40, 10] }];
  public chartLabels: Array<any> = ['Phizer', 'Sinopharm', 'Astrazeneca'];
  //graficoresult
  public chartDatasets2:Array<any> = [/* { data: [50, 40, 10] } */];
  public chartLabels2: Array<any> = ['Vivo', 'Muerto'];



  public listSexos: any[] = [
    { valor: 0, nombre: 'Masculino' },
    { valor: 1, nombre: 'Femenino' },
  ];
  public listCamaUCI: any[] = [
    { valor: 0, nombre: 'No' },
    { valor: 1, nombre: 'Si' },
  ];
  public listOxigeno: any[] = [
    { valor: 0, nombre: 'No' },
    { valor: 1, nombre: 'Si' },
  ];
  public listVentilacion: any[] = [
    { valor: 0, nombre: 'No' },
    { valor: 1, nombre: 'Si' },
  ];
  public listDosis1: any[] = [
    { valor: 0, nombre: 'No' },
    { valor: 1, nombre: 'Si' },
  ];
  public listDosis2: any[] = [
    { valor: 0, nombre: 'No' },
    { valor: 1, nombre: 'Si' },
  ];
  public listFabricante: any[] = [
    { valor: 0.35, nombre: 'Phizer' },
    { valor: 0.65, nombre: 'Sinopharm' },
    { valor: 1, nombre: 'Astrazeneca' },
  ];

  constructor(private fb: FormBuilder, private service: ServicesService) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      dni : [''],
      sexo: [''],
      edad: [''],
      flag_uci: [''],
      con_oxigeno: [''],
      con_ventilacion: [''],
      flag_vacuna1: [''],
      flag_vacuna2: [''],
      fabricante_dosis: [''],
    });
  }

  submitEnviar() {
    const data = {
      dni: this.form.value.dni,
      sexo: this.form.value.sexo,
      edad: this.form.value.edad  / 100,
      flag_uci: this.form.value.flag_uci,
      con_oxigeno: this.form.value.con_oxigeno,
      con_ventilacion: this.form.value.con_ventilacion,
      flag_vacuna1: this.form.value.flag_vacuna1,
      flag_vacuna2: this.form.value.flag_vacuna2,
      fabricante_dosis: this.form.value.fabricante_dosis,
    };
    console.log('[Form] data request: ', data);
    

    this.service.registerUser(data).subscribe((resp) => {
      this.isResultado = true;
      this.porcentaje = 0.8; /* resp.porcentaje */
      this.porcentajeTotal = 1 - this.porcentaje;

      //dibujar el grafico
      this.chartDatasets2 = [
        {data:[this.porcentaje, this.porcentajeTotal]}
      ]
    });
  }

  selectOpt(option: any) {
    console.log('option:: ', option);
    this.form.controls['fabricante_dosis'].setValue(option.value);
  }
}
