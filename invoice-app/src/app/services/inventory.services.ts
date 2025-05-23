import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class InventoryService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) { }

  createProduct(product: any) {
    return this.http.post(`${this.apiUrl}/products`, product);
  }

  getProduct(id: number) {
    return this.http.get(`${this.apiUrl}/products/${id}`);
  }
  getAllProducts(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/products`);
  }
}
