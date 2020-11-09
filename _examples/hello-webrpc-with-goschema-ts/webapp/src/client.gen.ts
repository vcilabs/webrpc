/* tslint:disable */
//   
// --
// This file has been generated by https://github.com/webrpc/webrpc using gen/typescript
// Do not edit by hand. Update your webrpc schema and re-generate.


//
// Types
//


export enum Kind {
}



export interface Empty {
}



export interface Page {
  num: number
}



export interface User {
  id: number
  username: string
  role: Kind
  meta: {[key: string]: any}
  internalID: number
}

export interface ExampleService {
  findUsers(args: FindUsersArgs, headers?: object): Promise<FindUsersReturn>
  getUser(args: GetUserArgs, headers?: object): Promise<GetUserReturn>
  ping(headers?: object): Promise<PingReturn>
}

export interface FindUsersArgs {
  q: string
}

export interface FindUsersReturn {
  page: Page
  user: Array<User>  
}
export interface GetUserArgs {
  userID: number
}

export interface GetUserReturn {
  user: User  
}
export interface PingArgs {
}

export interface PingReturn {
  bool: boolean  
}


  
//
// Client
//
export class ExampleService implements ExampleService {
  private hostname: string
  private fetch: Fetch
  private path = '/rpc/ExampleService/'

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = fetch
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }
  
  findUsers = (args: FindUsersArgs, headers?: object): Promise<FindUsersReturn> => {
    return this.fetch(
      this.url('FindUsers'),
      createHTTPRequest(args, headers)).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          page: <Page>(_data.page), 
          user: <Array<User>>(_data.user)
        }
      })
    })
  }
  
  getUser = (args: GetUserArgs, headers?: object): Promise<GetUserReturn> => {
    return this.fetch(
      this.url('GetUser'),
      createHTTPRequest(args, headers)).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          user: <User>(_data.user)
        }
      })
    })
  }
  
  ping = (headers?: object): Promise<PingReturn> => {
    return this.fetch(
      this.url('Ping'),
      createHTTPRequest({}, headers)
      ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          bool: <boolean>(_data.bool)
        }
      })
    })
  }
  
}

  
export interface WebRPCError extends Error {
  code: string
  msg: string
	status: number
}

const createHTTPRequest = (body: object = {}, headers: object = {}): object => {
  return {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify(body || {})
  }
}

const buildResponse = (res: Response): Promise<any> => {
  return res.text().then(text => {
    let data
    try {
      data = JSON.parse(text)
    } catch(err) {
      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status: res.status } as WebRPCError
    }
    if (!res.ok) {
      throw data // webrpc error response
    }
    return data
  })
}

export type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>
