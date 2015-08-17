<?php

namespace newsletters\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Support\Facades\Auth;

class ImportSubscribersRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize()
    {
        return Auth::check();
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array
     */
    public function rules()
    {
        return [
            'subscribers' => 'mimes:csv,xls,xlsx,txt|check_fields:subscribers,list_id',
            'list_id' => 'required'
        ];
    }
}
