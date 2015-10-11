<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 29.7.15
 * Time: 22:40
 */

namespace newsletters\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Support\Facades\Auth;

class StoreUserSettingsRequest extends FormRequest
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
            'name'       => 'required',
            'email'      => 'required|email',
            'password'   => 'sometimes|required',
            'aws_key'    => 'required',
            'aws_secret' => 'required',
            'aws_region' => 'required',
        ];
    }
}
