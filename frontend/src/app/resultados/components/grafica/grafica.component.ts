import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-grafica',
  templateUrl: './grafica.component.html',
  styleUrls: ['./grafica.component.scss']
})
export class GraficaComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

   // Grafico
   public chartType: string = 'pie';

   @Input() chartDatasets: Array<any> = [
     /* { data: [50, 40, 10]} */
   ];
 
   @Input() chartLabels: Array<any> = [/* 'Vacuna A', 'Vacuna B', 'Vacuna C' */];
 
   public chartColors: Array<any> = [
     {
       backgroundColor: ['#F7464A', '#46BFBD', '#FDB45C', '#949FB1', '#4D5360'],
       hoverBackgroundColor: ['#FF5A5E', '#5AD3D1', '#FFC870', '#A8B3C5', '#616774'],
       borderWidth: 5,
     }
   ];
 
   public chartOptions: any = {
     responsive: true
   };
   public chartClicked(e: any): void { }
   public chartHovered(e: any): void { }
   // Grafico

}
