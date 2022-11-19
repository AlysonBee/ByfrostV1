import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map,catchError, tap } from "rxjs/operators";
import { throwError } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ByfrostApiService {
  url: string = "http://localhost:4201/hello-indexer"
  urlPost: string = "http://localhost:4201/hello-indexer/request" 

  constructor(private http: HttpClient) { }

  getCode(labelToRequest: string) {
    let queryParams = new HttpParams()

    queryParams.append("token", labelToRequest)
    return this.http.get(`${this.url}?token=${labelToRequest}`,
        {
          params: queryParams
        }
      )
      .pipe(
        map(
          responseData => {
            return responseData
          }
        ),
        catchError(errorRes => {
          return throwError(errorRes);
      })
    );
  }


  requestBlock(labelToRequest: string) {
    let queryParams = new HttpParams()

    queryParams.append("token", labelToRequest)
    return this.http.post(this.urlPost,
      {
        "token": labelToRequest,
      },
    ).pipe(catchError(errorRes => {
        console.log("errorRes ", errorRes)
        return errorRes
      }), tap(responseData => {
        // console.log(responseData)
      })
    );
  }
}
