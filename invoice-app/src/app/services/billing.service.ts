import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class BillingService {

  private baseUrl = 'http://localhost:8081';

  constructor(private http: HttpClient) { }

  createInvoice(number: string): Observable<any> {
    return this.http.post(`${this.baseUrl}/invoices`, { number });
  }


  getInvoices(): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/invoices`);
  }


  getInvoice(invoiceId: number): Observable<any> {
    return this.http.get(`${this.baseUrl}/invoices/${invoiceId}`);
  }


  addItemToInvoice(invoiceId: number, productId: number, quantity: number): Observable<any> {
    return this.http.post(`${this.baseUrl}/invoices/${invoiceId}/items`, { product_id: productId, quantity });
  }


  printInvoice(invoiceId: number): Observable<any> {
    return this.http.post(`${this.baseUrl}/invoices/${invoiceId}/print`, {});
  }
}
