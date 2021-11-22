import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ServicesService {
  constructor(private http: HttpClient) {}

  public registerUser(data: any) {
    return this.http.post('http://localhost:8080/api/v1/appointments',data);
  }
}
