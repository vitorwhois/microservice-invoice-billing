import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class BillingService {
  private apiUrl = 'http://localhost:8081';

  constructor(private http: HttpClient) { }

  createInvoice(invoiceNumber: string) {
    return this.http.post(`${this.apiUrl}/invoices`, { number: invoiceNumber });
  }

  getInvoice(id: number) {
    return this.http.get(`${this.apiUrl}/invoices/${id}`);
  }

  addItemToInvoice(invoiceId: number, productId: number, quantity: number) {
    return this.http.post(`${this.apiUrl}/invoices/${invoiceId}/items`, {
      product_id: productId,
      quantity: quantity
    });
  }

  printInvoice(invoiceId: number) {
    return this.http.post(`${this.apiUrl}/invoices/${invoiceId}/print`, {});
  }
}
